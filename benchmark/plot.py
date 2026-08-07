import json
import os
import matplotlib.pyplot as plt
import matplotlib.ticker as ticker

def load_data(filepath):
    with open(filepath, 'r') as f:
        return json.load(f)

def plot_fpr_vs_elements(data, output_dir):
    fpr_data = data.get("fpr_vs_elements", [])
    if not fpr_data:
        print("No data found for FPR vs Elements experiment.")
        return

    elements = [row["inserted_elements"] for row in fpr_data]
    fill_ratios = [row["fill_ratio"] * 100 for row in fpr_data]
    empirical_fpr = [row["empirical_fpr"] for row in fpr_data]
    theoretical_fpr = [row["theoretical_fpr"] for row in fpr_data]

    # Use a clean, modern aesthetic style
    fig, ax1 = plt.subplots(figsize=(11, 7))

    # Set background styling
    fig.patch.set_facecolor('#fdfdfd')
    ax1.set_facecolor('#f7f7f7')

    # Primary axis: False Positive Rate vs. Fill Ratio (%)
    line1, = ax1.plot(fill_ratios, empirical_fpr, label='Empirical FPR', color='#1f77b4', linewidth=2.5, marker='o', markersize=4)
    line2, = ax1.plot(fill_ratios, theoretical_fpr, label='Theoretical FPR', color='#ff7f0e', linewidth=2, linestyle='--')
    
    # Detailed X and Y limits and scale increments
    ax1.set_xlim(0, 100)
    ax1.set_ylim(-0.02, 1.02)
    
    # Granular major/minor gridlines and ticks
    ax1.xaxis.set_major_locator(ticker.MultipleLocator(10))
    ax1.xaxis.set_minor_locator(ticker.MultipleLocator(2))
    ax1.yaxis.set_major_locator(ticker.MultipleLocator(0.1))
    ax1.yaxis.set_minor_locator(ticker.MultipleLocator(0.02))
    
    # Enable grids
    ax1.grid(True, which='major', linestyle='--', alpha=0.7, color='#aaaaaa')
    ax1.grid(True, which='minor', linestyle=':', alpha=0.4, color='#cccccc')
    ax1.minorticks_on()
    
    ax1.set_xlabel('Filter Saturation (Fill Ratio %)', fontsize=12, fontweight='bold', labelpad=10)
    ax1.set_ylabel('False Positive Rate', fontsize=12, fontweight='bold', labelpad=10)
    ax1.tick_params(axis='both', which='major', labelsize=10)
    ax1.tick_params(axis='both', which='minor', labelsize=8, colors='#777777')

    # Secondary X-axis on top to show Number of Elements (n)
    ax2 = ax1.twiny()
    ax2.set_xlim(ax1.get_xlim())
    # Interpolate element counts for ticks
    num_ticks = 11
    tick_indices = [int(i * (len(fill_ratios) - 1) / (num_ticks - 1)) for i in range(num_ticks)]
    ax2.set_xticks([fill_ratios[idx] for idx in tick_indices])
    ax2.set_xticklabels([f"{elements[idx]:,}" for idx in tick_indices])
    ax2.set_xlabel('Number of Elements Filled ($n$)', fontsize=11, color='#555555', labelpad=12)
    ax2.tick_params(axis='x', colors='#555555', labelsize=9)

    plt.title('Bloom Filter False Positive Rate vs. Filter Saturation (m = 100,000, k = 4)', 
              fontsize=14, fontweight='bold', pad=22, color='#333333')
    
    # Legend
    lines = [line1, line2]
    labels = [l.get_label() for l in lines]
    ax1.legend(lines, labels, loc='upper left', frameon=True, facecolor='white', edgecolor='#e0e0e0', fontsize=11)

    # Layout adjustment and saving
    plt.tight_layout()
    output_path = os.path.join(output_dir, 'fpr_vs_elements.png')
    plt.savefig(output_path, dpi=300, facecolor=fig.get_facecolor(), edgecolor='none')
    plt.close()
    print(f"Generated chart: {output_path}")

def plot_fpr_vs_k(data, output_dir):
    k_data = data.get("fpr_vs_k", [])
    if not k_data:
        print("No data found for FPR vs K experiment.")
        return

    k_values = [row["k"] for row in k_data]
    empirical_fpr = [row["empirical_fpr"] for row in k_data]
    theoretical_fpr = [row["theoretical_fpr"] for row in k_data]

    # Create figure
    fig, ax = plt.subplots(figsize=(11, 7))

    # Styling
    fig.patch.set_facecolor('#fdfdfd')
    ax.set_facecolor('#f7f7f7')

    # Plot lines
    ax.plot(k_values, empirical_fpr, label='Empirical FPR', color='#2ca02c', linewidth=2.5, marker='s', markersize=5)
    ax.plot(k_values, theoretical_fpr, label='Theoretical FPR', color='#d62728', linewidth=2, linestyle='--')

    # Granular major/minor gridlines and ticks
    ax.set_xlim(0.5, 16.5)
    ax.set_ylim(-0.02, max(max(empirical_fpr), max(theoretical_fpr)) + 0.05)
    
    ax.xaxis.set_major_locator(ticker.MultipleLocator(1))
    ax.yaxis.set_major_locator(ticker.MultipleLocator(0.05))
    ax.yaxis.set_minor_locator(ticker.MultipleLocator(0.01))
    
    # Enable grids
    ax.grid(True, which='major', linestyle='--', alpha=0.7, color='#aaaaaa')
    ax.grid(True, which='minor', linestyle=':', alpha=0.4, color='#cccccc')
    ax.minorticks_on()

    # Find and highlight optimal k
    min_theoretical_idx = theoretical_fpr.index(min(theoretical_fpr))
    opt_k = k_values[min_theoretical_idx]
    opt_val = theoretical_fpr[min_theoretical_idx]
    ax.axvline(x=opt_k, color='#7f7f7f', linestyle=':', alpha=0.8, linewidth=1.5)
    ax.annotate(f'Optimal k = {opt_k}\n(FPR = {opt_val:.4f})', 
                xy=(opt_k, opt_val), 
                xytext=(opt_k + 1.2, opt_val + 0.05),
                arrowprops=dict(facecolor='#555555', shrink=0.08, width=1.5, headwidth=6, headlength=6),
                fontsize=10, fontweight='semibold', color='#333333',
                bbox=dict(boxstyle="round,pad=0.3", fc="white", ec="#e0e0e0", alpha=0.9))

    # Axis properties
    ax.set_xlabel('Number of Hash Functions ($k$)', fontsize=12, fontweight='bold', labelpad=10)
    ax.set_ylabel('False Positive Rate', fontsize=12, fontweight='bold', labelpad=10)
    ax.tick_params(axis='both', which='major', labelsize=10)
    ax.tick_params(axis='both', which='minor', labelsize=8, colors='#777777')
    
    plt.title('Bloom Filter False Positive Rate vs. Number of Hash Functions k (m = 100,000, n = 20,000)', 
              fontsize=14, fontweight='bold', pad=22, color='#333333')
    
    ax.legend(loc='upper right', frameon=True, facecolor='white', edgecolor='#e0e0e0', fontsize=11)

    plt.tight_layout()
    output_path = os.path.join(output_dir, 'fpr_vs_k.png')
    plt.savefig(output_path, dpi=300, facecolor=fig.get_facecolor(), edgecolor='none')
    plt.close()
    print(f"Generated chart: {output_path}")

def main():
    json_path = 'benchmark_results.json'
    output_dir = 'benchmark'
    
    if not os.path.exists(output_dir):
        os.makedirs(output_dir)

    if not os.path.exists(json_path):
        print(f"Error: {json_path} not found. Please run the benchmark tool first.")
        return

    print("Loading benchmark results...")
    data = load_data(json_path)

    print("Generating charts...")
    plot_fpr_vs_elements(data, output_dir)
    plot_fpr_vs_k(data, output_dir)
    print("Done!")

if __name__ == '__main__':
    main()
